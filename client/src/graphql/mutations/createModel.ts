import { gql } from '@apollo/client'

export default gql`
    mutation CreateModel($name: String!, $fields: [ModelFieldInput!]!) {
        createModel(name: $name, fields: $fields) {
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
