require 'formula'

class Muxy < Formula
  homepage "https://github.com/packer-community/packer-windows-plugins"
  version "0.0.1"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_amd64.zip"
    sha256 'ba29e1b5c2666a91f37198ebef4397d497d262c3a4530eeed96d206824481634'
  else
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_386.zip"
    sha256 '9dad1e255359738c74d359b3e5e0580ad79f296b4cd6d7cfb350d5055b00fd87'
  end

  depends_on :arch => :intel

  def install
    bin.install Dir['*']
  end

  test do

  end
end
