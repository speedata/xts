desc "Create a GitHub release (needs a git tag)"
task "release" do
	sh "goreleaser release --clean"
	sh "xcrun notarytool submit -p notary dist/xts_macos_intel.zip --wait"
	sh "xcrun notarytool submit -p notary dist/xts_macos_arm64.zip --wait"
end

desc "Update markdown documentation"
task "doc" => [:xtshelper] do
	sh "#{INSTALDIR}/bin/xtshelper doc"
	puts "done"
end
