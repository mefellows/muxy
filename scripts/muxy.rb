require 'formula'

class Muxy < Formula
  homepage "https://github.com/packer-community/packer-windows-plugins"
  version "0.0.1"

  if Hardware.is_64_bit?
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_amd64.zip"
    sha1 'f9003414fadb5fb483a299866018eb6b7772d329'
  else
    url "https://github.com/mefellows/muxy/releases/download/v#{version}/darwin_386.zip"
    sha1 'a6bfa0bf2d9b012503a375735987aeed9f12cb6c'
  end

  depends_on :arch => :intel

  def install
    bin.install Dir['*']
  end

  test do
    
  end
end
