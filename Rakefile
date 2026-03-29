require "pathname"

# Get version from git tag (e.g., "v1.0.0" or "v1.0.0-3-g1a2b3c4")
def git_version
  version = `git describe --tags --always --match 'v*' 2>/dev/null`.strip
  version.empty? ? "dev" : version.sub(/^v/, "")
end

INSTALDIR = Pathname.new(__dir__).freeze

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
task :qa do
	sh "xts compare #{INSTALDIR}/qa"
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
