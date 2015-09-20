require 'formula'

class Muxy < Formula
  homepage "https://github.com/packer-community/packer-windows-plugins"
  version "0.0.1"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_amd64.zip"
    sha1 'b8f75222acc0ff04c6ead6527b6e3e1d2a3c8a57'
  else
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_386.zip"
    sha1 '47821cb8b2b9e112aa820aaa2ec7396053ea9582'
  end

  depends_on :arch => :intel

  def install
    bin.install Dir['*']
  end

  test do
    
  end
end
