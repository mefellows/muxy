require 'formula'

class Muxy < Formula
  homepage "https://github.com/packer-community/packer-windows-plugins"
  version "0.0.1"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_amd64.zip"
    sha1 'a2e54979a03a0cc4b30527c6c7d61e97f23cbf11'
  else
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_386.zip"
    sha1 '9d942be55cf5af6c17cc68f9beee6a60538105f1'
  end

  depends_on :arch => :intel

  def install
    bin.install Dir['*']
  end

  test do
    
  end
end
