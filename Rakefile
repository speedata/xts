require "pathname"
@versions = {}

File.read("version").each_line do |line|
	product,versionnumber = line.chomp.split(/=/)
	@versions[product]=versionnumber
end

installdir = Pathname.new(__FILE__).join("..")

@xts_version = @versions['xts_version']

desc "Show rake description"
task :default do
	puts
	puts "Run 'rake -T' for a list of tasks."
	puts
	puts "1: Use 'rake build' to build the 'xtss' binary. That should be\n   the starting point."
	puts
end

task :xtshelper  do
	sh "go build -ldflags \"-X main.version=#{@xts_version} -X main.basedir=#{installdir} \" -o bin/xtshelper github.com/speedata/xts/helper"
end


desc "Compile and install necessary software"
task :build do
	sh "go build -ldflags \"-X github.com/speedata/xts/core.Version=#{@xts_version}\" -o bin/xts github.com/speedata/xts/xts"
end

task "doc" => [:xtshelper] do
	sh "#{installdir}/bin/xtshelper doc"
	puts "done"
end

desc "Create a GitHub release (needs a git tag)"
task "release" do
	sh "goreleaser release --clean"
	sh "xcrun notarytool submit -p notary dist/xts_macos_intel.zip --wait"
	sh "xcrun notarytool submit -p notary dist/xts_macos_arm64.zip --wait"
end

task :schema => [:xtshelper] do
	sh "bin/xtshelper genschema"
end
