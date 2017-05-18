require 'formula'

class Muxy < Formula
  homepage "https://github.com/packer-community/packer-windows-plugins"
  version "0.0.4"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_amd64.zip"
    sha256 '6d4aa05dfd0d94c98e93b6306bf660228e7254108c9f970e84932be8087d93b5'
  else
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_386.zip"
    sha256 '51c6b34846cb1a2913e36af9b4eefe29d71b1d5efe053a720a500ec2f52a0378'
  end

  depends_on :arch => :intel

  def install
    bin.install Dir['*']
  end

  test do

  end
end
