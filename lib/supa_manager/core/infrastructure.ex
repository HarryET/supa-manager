defmodule SupaManager.Core.Infrastructure do
  @spec aws_region(String.t()) :: String.t()
  def aws_region(input) do
    case String.downcase(input) do
      "west us (north california)" -> "us-west-1"
      "east us (north virginia)" -> "us-east-1"
      "canada (central)" -> "ca-central-1"
      "west eu (ireland)" -> "eu-west-1"
      "west eu (london)" -> "eu-west-2"
      "north eu" -> "eu-north-1"
      "central eu (frankfurt)" -> "eu-central-1"
      "south asia (mumbai)" -> "ap-south-1"
      "southeast asia (singapore)" -> "ap-southeast-1"
      "northeast asia (tokyo)" -> "ap-northeast-1"
      "northeast asia (seoul)" -> "ap-northeast-2"
      "oceania (sydney)" -> "ap-southeast-2"
      "south america (são paulo)" -> "sa-east-1"
      _ -> "unknown"
    end
  end
end
