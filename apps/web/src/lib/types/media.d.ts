import { PageDTO, MediaOverviewDTO } from "../../dto"

export type FetchMedia = (params: URLSearchParams) => Promise<PageDTO<MediaOverviewDTO>>