import { ScanPathData, GenerateThumbnailData, type RefreshMetadata, type RefreshLibraryMetadata, type GenerateChaptersData } from "../../dto"

export type JobData =
  | ScanPathData
  | GenerateThumbnailData
  | RefreshMetadata
  | RefreshLibraryMetadata
  | GenerateChaptersData