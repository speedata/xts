require "pathname"

# Get version from git tag (e.g., "v1.0.0" or "v1.0.0-3-g1a2b3c4")
def git_version
  version = `git describe --tags --always --match 'v*' 2>/dev/null`.strip
  version.empty? ? "dev" : version.sub(/^v/, "")
end

INSTALDIR = Pathname.new(__dir__).freeze

# The binary `rake build` produces. QA tasks must invoke this path and never a
# bare `xts`: a bare name resolves through PATH to whatever was last installed
# into $GOBIN, so the QA would silently compare against a stale build instead
# of the working tree.
XTSBIN = INSTALDIR.join("bin", "xts").freeze

@xts_version = git_version

desc "Show rake description"
task :default do
	puts
	puts "Run 'rake -T' for a list of tasks."
	puts
	puts "1: Use 'rake build' to build the 'xts' binary. That should be\n   the starting point."
	puts
end

task :xtshelper  do
	sh "go build -ldflags \"-X main.version=#{@xts_version} -X main.basedir=#{INSTALDIR} \" -o bin/xtshelper github.com/speedata/xts/helper"
end

desc "Create markdown reference and the changelog page"
task :doc => [:xtshelper] do
	sh "bin/xtshelper doc"
end

desc "Build the manual with Hugo and check its internal links with htmltest"
task :doccheck => [:doc] do
	# The manual is deployed under https://doc.speedata.de/xts/, so the check
	# builds it into the subdirectory xts/ of an otherwise empty web root and
	# lets htmltest resolve links against that root. A root-relative path such
	# as /manual/img/x.png then fails here as it does on the server.
	webroot = INSTALDIR.join("build", "manual").to_s
	FileUtils.rm_rf(webroot)
	sh "cd doc/manual && hugo --quiet --baseURL http://localhost/xts/ -d #{webroot}/xts"
	sh "htmltest -c doc/manual/.htmltest.yml #{webroot}"
end

desc "Render the manual images from doc/manual/images/*/layout.xml: rake docimages[name]"
task :docimages, [:name] => [:build] do |t, args|
	# Every subdirectory of doc/manual/images holds the layout (and the
	# optional data.xml and overlay files) of one image in the manual. The
	# first page of its PDF becomes content/manual/img/<name>.png, so the
	# page format of the layout is the image size. See doc/manual/images/Readme.md.
	srcdir = INSTALDIR.join("doc", "manual", "images")
	outdir = INSTALDIR.join("doc", "manual", "content", "manual", "img")
	dirs = args[:name] ? [srcdir.join(args[:name])] : srcdir.children.sort
	dirs.select { |d| d.join("layout.xml").exist? }.each do |dir|
		data = dir.join("data.xml").exist? ? "" : "--dummy"
		sh "cd #{dir} && #{XTSBIN} --suppressinfo --quiet #{data}"
		sh "pdftoppm -r 200 -png -singlefile #{dir.join('xts.pdf')} #{outdir.join(dir.basename)}"
	end
end

desc "Check that doc/changelog.xml is ready for a release: rake changelog[v0.1.0]"
task :changelog, [:version] => [:xtshelper] do |t, args|
	version = args[:version] || suggest_next_version
	sh "bin/xtshelper changelog check #{version}"
	puts "Release notes for #{version}:"
	puts
	sh "bin/xtshelper changelog notes #{version}"
end

desc "Build the 'xts' binary"
task :build do
	sh "go build -ldflags \"-s -w -X github.com/speedata/xts/core.Version=#{@xts_version}\" -o bin/xts github.com/speedata/xts/xts"
end

desc "Install 'xts' into $GOBIN"
task :install do
	sh "go install -ldflags \"-s -w -X github.com/speedata/xts/core.Version=#{@xts_version}\" github.com/speedata/xts/xts"
end

desc "Create the schema files"
task :schema => [:xtshelper] do
	sh "bin/xtshelper genschema"
end


desc "Run quality assurance"
task :qa => [:build] do
	sh XTSBIN.to_s, "compare", INSTALDIR.join("qa").to_s
end

desc "Clean QA intermediate files"
task :cleanqa do
	FileUtils.rm Dir.glob("qa/**/pagediff-*.png")
	FileUtils.rm Dir.glob("qa/**/reference-*.png")
	FileUtils.rm Dir.glob("qa/**/source-*.png")
	FileUtils.rm Dir.glob("qa/**/xts-aux.xml")
	FileUtils.rm Dir.glob("qa/**/xts-protocol.xml")
	FileUtils.rm Dir.glob("qa/**/xts.pdf")
end
