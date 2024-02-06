import { ApolloProvider } from '@apollo/client'
import { useState } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { ModelsProvider } from './context/ModelsContext'
import { useApolloClient } from './hooks/useApolloClient'
import Home from './pages/Home/Home'

const WrappedApp = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path={'/'} element={<Home />} />
            </Routes>
        </BrowserRouter>
    )
}

const App = () => {
    const [error, setError] = useState<Error | null>(null)
    const apolloClient = useApolloClient('http://localhost:8080/next-public', setError)

    if (!apolloClient) {
        return null
    }

    console.log('error', error)
    // if (error) {
    //     return <Box>{error.message}</Box>
    // }

    return (
        <ApolloProvider client={apolloClient}>
            <ModelsProvider>
                <WrappedApp />
            </ModelsProvider>
        </ApolloProvider>
    )
}

export default App
