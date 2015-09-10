require 'formula'

class Muxy < Formula
  homepage "https://github.com/packer-community/packer-windows-plugins"
  version "0.0.1"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_amd64.zip"
    sha1 'edb5b7abaf6fd5d03acfdbe78564d3f60a4ddd7a'
  else
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_386.zip"
    sha1 'fa0b81d2d59d750efe8ffc862fddd909da505a95'
  end

  depends_on :arch => :intel

  def install
    bin.install Dir['*']
  end

  test do
    
  end
end
