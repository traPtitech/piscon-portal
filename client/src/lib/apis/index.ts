import { Apis, Configuration } from './generated'

const apis = new Apis(new Configuration({ basePath: '/api' }))

export default apis
export * from './generated'
