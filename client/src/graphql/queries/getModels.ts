import { gql } from '@apollo/client'

export default gql`
    query GetModels {
        getModels {
            name
            fields {
                name
                type
                isNullable
                subFields {
                    name
                    type
                    isNullable
                    subFields {
                        name
                        type
                        isNullable
                    }
                }
            }
        }
    }
`
