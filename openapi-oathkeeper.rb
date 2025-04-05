# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class OpenapiOathkeeper < Formula
  desc "This project aims to automating the generation of Ory Oathkeeper rules from an OpenAPI 3 contract and save a lot of time and effort, especially for larger projects with many endpoints or many services."
  homepage "https://github.com/cerberauth/openapi-oathkeeper"
  version "0.7.11"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/cerberauth/openapi-oathkeeper/releases/download/v0.7.11/openapi-oathkeeper_Darwin_x86_64.tar.gz"
      sha256 "4beb8df051a67ed2130e638720e8ade8e1e2385d46ea059be1abbc869f0fb75b"

      def install
        bin.install "openapi-oathkeeper"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/cerberauth/openapi-oathkeeper/releases/download/v0.7.11/openapi-oathkeeper_Darwin_arm64.tar.gz"
      sha256 "22b430758ded37900d58a266f07a10f4f4b0822c5738edbe0a87d42a2a1b63b2"

      def install
        bin.install "openapi-oathkeeper"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/cerberauth/openapi-oathkeeper/releases/download/v0.7.11/openapi-oathkeeper_Linux_x86_64.tar.gz"
        sha256 "79fadd8b19f8ffe9d7fb3dc803f0be28c1fb121aa0f4fa4a2a4820b55d0e1c2a"

        def install
          bin.install "openapi-oathkeeper"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/cerberauth/openapi-oathkeeper/releases/download/v0.7.11/openapi-oathkeeper_Linux_arm64.tar.gz"
        sha256 "3015d8b1d4749146bfeb3a4ba66a487bf6cc94ce27fb2524d054fdf4e4320795"

        def install
          bin.install "openapi-oathkeeper"
        end
      end
    end
  end

  test do
    system "#{bin}/openapi-oathkeeper help"
  end
end
