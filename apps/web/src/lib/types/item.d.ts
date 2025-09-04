export type Item = {
  id: string
  name: string
}

export type ItemUrlFn = (id: string, name: string) => string

export type FetchItems = () => Promise<Item[]>

export type RemoveItem = (itemId: string) => Promise

export type AddItem = (itemId: string) => Promise

export type CreateItem = (name: string) => Promise<Item>

export type SearchItems = (name: string) => Promse<Item[]>

export type SearchItemsWithParams = (params: URLSearchParams) => Promise<Item[]>

