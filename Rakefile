require "pathname"
@versions = {}

if File.file?("version") then
	File.read("version").each_line do |line|
		product,versionnumber = line.chomp.split(/=/)
		@versions[product]=versionnumber
	end
else
	@versions["xts_version"]="1.0.0"
end

INSTALDIR = Pathname.new(__dir__).freeze

@xts_version = @versions['xts_version']

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


desc "Compile and install necessary software"
task :build do
	sh "go build -ldflags \"-X github.com/speedata/xts/core.Version=#{@xts_version}\" -o bin/xts github.com/speedata/xts/xts"
end

desc "Create the schema files"
task :schema => [:xtshelper] do
	sh "bin/xtshelper genschema"
end

desc "Update the version information from the latest git tag"
task :updateversion do
	sh "git describe| sed s,v,xts_version=, > version"
end


desc "Run quality assurance"
task :qa do
	sh "#{INSTALDIR}/bin/xts compare #{INSTALDIR}/qa"
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
