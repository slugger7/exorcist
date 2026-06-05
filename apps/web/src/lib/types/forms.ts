export type Validator = (value: any, state: any) => string | undefined

export type ValidationHandler = (
  { value, state, validators }:
    { value: any, state: any, validators: Validator[] | undefined }) => string[] | undefined

export type Touched<T> = Partial<{
  [K in keyof T]: T[K] extends object ? Touched<T[K]> : boolean
}>

export type Errors<T> = Partial<{
  [K in keyof T]: T[K] extends object ? Errors<T[K]> : string[] | undefined
}>

export type Validators<T> = Partial<{
  [K in keyof T]: T[K] extends object ? Validators<T[K]> : Validator[]
}>

export type GetErrorFn<T> = (key: keyof T) => string[] | undefined
