import { Apis, Configuration } from './generated'

const apis = new Apis(new Configuration({}))

export default apis
export * from './generated'
