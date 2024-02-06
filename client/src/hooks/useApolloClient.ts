import { ApolloClient, HttpLink, InMemoryCache, NormalizedCacheObject, from } from '@apollo/client'
import { onError } from '@apollo/client/link/error'
import { useCallback, useEffect, useState } from 'react'

export const useApolloClient = (graphQLServerURL: string, setError: (error: Error) => void) => {
    const [apolloClient, setApolloClient] = useState<ApolloClient<NormalizedCacheObject>>()

    const errorLink = useCallback(
        () =>
            onError((e) => {
                console.log('Error: ', e.graphQLErrors, e.networkError)
                const errors = {
                    graphQLErrors: e.graphQLErrors?.map((error) => error.message),
                    networkError: e.networkError?.message,
                }
                const error = new Error(JSON.stringify(errors))
                setError(error)
                setApolloClient(undefined)
            }),
        [setError],
    )

    useEffect(() => {
        console.log('Creating Apollo Client')
        const httpLink = new HttpLink({ uri: graphQLServerURL })

        const client = new ApolloClient({
            cache: new InMemoryCache(),
            link: from([errorLink(), httpLink]),
        })
        setApolloClient(client)
    }, [errorLink, graphQLServerURL])

    return apolloClient
}
