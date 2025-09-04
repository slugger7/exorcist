import type { WSTopicAllValues } from "../../dto"
import type { PageDTO } from "../../dto";

export * from "../../dto"

// https://www.typescriptlang.org/docs/handbook/jsdoc-supported-types.html
export * from "./job"
export * from "./item"
export * from "./media"

export type LoginResult = {
  userId: string;
  username: string;
}

export type WSTopicMap<T> = {
  [key in WSTopicAllValues]?: (data: T) => void
}

export type WSMessage<T> = {
  topic: WSTopicAllValues
  data: T
}

export type Ordinal = {
  text: string
  value: string
}

export type WSTopicMapView<T> = {
  [key in WSTopicAllValues]?: (mediaPage: PageDTO<T>, newItems: T[], data: T) => [PageDTO<T>, T[]]
}