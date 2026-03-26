# Distribution tasks for xts
# Cross-compilation, code signing, notarization, and release

require 'digest'

PLATFORMS = [
  { os: "linux",   arch: "amd64", ext: "", name: "linux" },
  { os: "linux",   arch: "arm64", ext: "", name: "linux" },
  { os: "darwin",  arch: "amd64", ext: "", name: "macos" },
  { os: "darwin",  arch: "arm64", ext: "", name: "macos" },
  { os: "windows", arch: "amd64", ext: ".exe", name: "windows" },
]

# macOS signing identity (from: security find-identity -v -p codesigning)
SIGNING_IDENTITY = "Developer ID Application: PATRICK GUENTER GUNDLACH (3Y98DLKYBJ)"

# Notarization keychain profile (created with: xcrun notarytool store-credentials "notary")
NOTARY_PROFILE = "notary"

DIST_DIR = "dist"
HOMEBREW_TAP = File.expand_path("~/work/software/bag-singlerepos/homebrew-tap")

desc "Build binaries for all platforms"
task :crossbuild => :schema do
  FileUtils.mkdir_p(DIST_DIR)
  ldflags = "-s -w -X github.com/speedata/xts/core.Version=#{@xts_version}"

  PLATFORMS.each do |p|
    output = "#{DIST_DIR}/xts-#{p[:name]}-#{p[:arch]}#{p[:ext]}"
    puts "Building #{output}..."
    sh "GOOS=#{p[:os]} GOARCH=#{p[:arch]} CGO_ENABLED=0 go build -ldflags '#{ldflags}' -o #{output} github.com/speedata/xts/xts"
  end
end

desc "Sign macOS binaries"
task :sign => :crossbuild do
  darwin_binaries.each do |bin|
    puts "Signing #{bin}..."
    sh %Q{codesign --force --options runtime --sign "#{SIGNING_IDENTITY}" --timestamp "#{bin}"}
    sh %Q{codesign --verify --verbose "#{bin}"}
  end
end

desc "Create distribution archives"
task :archives => :crossbuild do
  schema_files = Dir["schema/*schema*"]

  PLATFORMS.each do |p|
    binary = "#{DIST_DIR}/xts-#{p[:name]}-#{p[:arch]}#{p[:ext]}"
    basename = "xts-#{p[:name]}-#{p[:arch]}"
    archive_dir = "#{DIST_DIR}/archive-tmp/#{basename}"
    target_binary = "xts#{p[:ext]}"

    # Create temp directory with all files
    FileUtils.rm_rf(archive_dir)
    FileUtils.mkdir_p("#{archive_dir}/schema")
    FileUtils.cp(binary, "#{archive_dir}/#{target_binary}")
    schema_files.each { |f| FileUtils.cp(f, "#{archive_dir}/schema/") }

    if p[:os] == "windows" || p[:os] == "darwin"
      sh %Q{cd "#{DIST_DIR}/archive-tmp" && zip -r "../#{basename}.zip" "#{basename}"}
    else
      sh %Q{tar -czf "#{DIST_DIR}/#{basename}.tar.gz" -C "#{DIST_DIR}/archive-tmp" "#{basename}"}
    end
  end

  # Clean up temp directory
  FileUtils.rm_rf("#{DIST_DIR}/archive-tmp")
end

desc "Notarize macOS distribution archives"
task :notarize => :archives do
  darwin_archives.each do |archive|
    puts "Notarizing #{archive}..."
    sh %Q{xcrun notarytool submit "#{archive}" --keychain-profile "#{NOTARY_PROFILE}" --wait}
  end
  puts "Notarization complete!"
end

desc "Create signed and notarized distribution"
task :dist do
  Rake::Task[:crossbuild].invoke
  if RUBY_PLATFORM =~ /darwin/
    Rake::Task[:sign].reenable
    Rake::Task[:sign].invoke
  end
  Rake::Task[:archives].reenable
  Rake::Task[:archives].invoke
  if RUBY_PLATFORM =~ /darwin/
    Rake::Task[:notarize].reenable
    Rake::Task[:notarize].invoke
  end
end

desc "Create GitHub release with all binaries"
task :release do
  # Working directory must be clean
  unless system("git diff --quiet && git diff --cached --quiet")
    abort "Working directory is not clean. Commit or stash your changes first."
  end

  # Check if HEAD has a tag
  existing_tag = `git tag --points-at HEAD`.strip.split("\n").first

  if existing_tag
    tag = existing_tag
    puts "Using existing tag: #{tag}"
  else
    suggested = suggest_next_version
    print "Release tag [#{suggested}]: "
    input = STDIN.gets.chomp
    tag = input.empty? ? suggested : input
    sh "git tag -a #{tag} -m 'Release #{tag}'"
  end

  sh "git push origin #{tag}"

  # Use the tag directly as version
  @xts_version = tag.sub(/^v/, "")

  # Clean old dist files, then rebuild
  FileUtils.rm_rf(DIST_DIR)
  Rake::Task[:crossbuild].reenable
  Rake::Task[:dist].reenable
  Rake::Task[:dist].invoke

  # Find all archives
  archives = Dir["#{DIST_DIR}/*.tar.gz"] + Dir["#{DIST_DIR}/*.zip"]

  if archives.empty?
    puts "No archives found in #{DIST_DIR}/"
    exit 1
  end

  puts "Creating release #{tag} with: #{archives.join(', ')}"
  sh %Q{gh release create "#{tag}" #{archives.join(' ')} --title "#{tag}" --generate-notes}

  # Update Homebrew formula
  Rake::Task[:update_formula].invoke(tag.sub(/^v/, ""))
end

desc "Update Homebrew formula with SHA256 from local archives"
task :update_formula, [:version] do |t, args|
  version = args[:version] || @xts_version
  formula_path = "#{HOMEBREW_TAP}/Formula/xts.rb"

  unless File.exist?(formula_path)
    puts "Homebrew formula not found at #{formula_path} — skipping"
    next
  end

  puts "Updating Homebrew formula for version #{version}..."

  # Compute SHA256 from local archives (macOS = .zip, Linux = .tar.gz)
  shas = {}
  {"macos-arm64" => "zip", "macos-amd64" => "zip", "linux-arm64" => "tar.gz", "linux-amd64" => "tar.gz"}.each do |platform, ext|
    archive = "#{DIST_DIR}/xts-#{platform}.#{ext}"
    if File.exist?(archive)
      shas[platform] = Digest::SHA256.file(archive).hexdigest
      puts "  #{platform}: #{shas[platform]}"
    else
      puts "  Warning: #{archive} not found"
    end
  end

  # Update formula
  content = File.read(formula_path)
  content.gsub!(/version ".*"/, %Q{version "#{version}"})

  shas.each do |platform, sha|
    case platform
    when "macos-arm64"
      content.gsub!(/(xts-macos-arm64\.zip"\n\s+sha256 ")[a-f0-9]+"/, "\\1#{sha}\"")
    when "macos-amd64"
      content.gsub!(/(xts-macos-amd64\.zip"\n\s+sha256 ")[a-f0-9]+"/, "\\1#{sha}\"")
    when "linux-arm64"
      content.gsub!(/(xts-linux-arm64\.tar\.gz"\n\s+sha256 ")[a-f0-9]+"/, "\\1#{sha}\"")
    when "linux-amd64"
      content.gsub!(/(xts-linux-amd64\.tar\.gz"\n\s+sha256 ")[a-f0-9]+"/, "\\1#{sha}\"")
    end
  end

  File.write(formula_path, content)
  puts "Updated #{formula_path}"
  puts "Don't forget to commit and push the homebrew-tap:"
  puts "  cd #{HOMEBREW_TAP} && git add Formula/xts.rb && git commit -m 'Update xts formula for version #{version}' && git push"
end

desc "Clean distribution artifacts"
task :distclean do
  FileUtils.rm_rf(DIST_DIR)
end

def darwin_binaries
  Dir["#{DIST_DIR}/xts-macos-*"].reject { |f| f.end_with?(".zip", ".tar.gz") }
end

def darwin_archives
  Dir["#{DIST_DIR}/xts-macos-*.zip"]
end

def suggest_next_version
  last = `git describe --tags --abbrev=0 2>/dev/null`.strip
  return "v0.1.0" if last.empty?

  if last =~ /^v?(\d+)\.(\d+)\.(\d+)$/
    "v#{$1}.#{$2}.#{$3.to_i + 1}"
  else
    last
  end
end
