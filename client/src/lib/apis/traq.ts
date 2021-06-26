import { Apis, Configuration } from '@traptitech/traq'
const BASE_PATH = '/api/v3'

const traqApis = new Apis(
  new Configuration({
    basePath: BASE_PATH
  })
)

export default traqApis
