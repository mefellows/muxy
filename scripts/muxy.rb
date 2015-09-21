require 'formula'

class Muxy < Formula
  homepage "https://github.com/packer-community/packer-windows-plugins"
  version "0.0.1"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_amd64.zip"
    sha1 '252001435b5aed95a4ec833e27fd422c53b08c42'
  else
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_386.zip"
    sha1 '805d2ae25ecabdd07dafaad10179499e7363a413'
  end

  depends_on :arch => :intel

  def install
    bin.install Dir['*']
  end

  test do

  end
end
