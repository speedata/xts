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

desc "Create markdown reference"
task :doc => [:xtshelper] do
	sh "bin/xtshelper doc"
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
