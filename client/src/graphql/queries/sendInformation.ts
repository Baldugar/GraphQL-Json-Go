import { gql } from '@apollo/client'

export default gql`
    query SendInformaton($info: Map!, $modelName: String!) {
        sendInformaton(info: $info, modelName: $modelName)
    }
`
