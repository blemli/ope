cask "ope" do
  version :latest
  sha256 :no_check

  url "https://github.com/blemli/ope/releases/latest/download/Ope.app.zip"
  name "Ope"
  desc "Open files and folders from the browser via ope:// URLs"
  homepage "https://github.com/blemli/ope"

  app "Ope.app"

  postflight do
    system_command "#{appdir}/Ope.app/Contents/MacOS/ope", args: ["install"]
  end

  uninstall_postflight do
    system_command "#{appdir}/Ope.app/Contents/MacOS/ope", args: ["uninstall"]
  end
end
