import { useLazyQuery, useMutation, useQuery } from '@apollo/client'
import { createContext, useState } from 'react'
import createModel from '../graphql/mutations/createModel'
import getModels from '../graphql/queries/getModels'
import sendInformation from '../graphql/queries/sendInformation'
import { Model, Mutation, MutationcreateModelArgs, Query, QuerysendInformatonArgs } from '../graphql/types'

export type ModelsContextType = {
    models: Model[]
    addModel: (model: Model) => void
    sendInformation: (information: { [key: string]: any }, modelName: string) => void
}

export const ModelsContext = createContext<ModelsContextType>({
    addModel: () => {
        console.error('No ModelsProvider')
    },
    models: [],
    sendInformation: () => {
        console.error('No ModelsProvider')
    },
})

const localStorageKey = 'models'

export const ModelsProvider = ({ children }: { children: React.ReactNode }) => {
    const [models, setModels] = useState<Model[]>(() => {
        const data = localStorage.getItem(localStorageKey)
        if (data) {
            return JSON.parse(data)
        }
        return []
    })

    useQuery<Query>(getModels, {
        skip: models.length > 0,
        onCompleted: (data) => {
            if (data.getModels) {
                setModels(data.getModels)
                localStorage.setItem(localStorageKey, JSON.stringify(data.getModels))
            }
        },
    })

    const [sendInformationQuery] = useLazyQuery<Query, QuerysendInformatonArgs>(sendInformation, {
        onCompleted: (data) => {
            if (data.sendInformaton === true) {
                alert('Information sent successfully')
            } else {
                alert('Failed to send information')
            }
        },
    })

    const handleSendInformation = (information: { [key: string]: any }, modelName: string) => {
        sendInformationQuery({
            variables: {
                info: information,
                modelName,
            },
        })
    }

    const [createModelMutation] = useMutation<Mutation, MutationcreateModelArgs>(createModel, {
        onCompleted: (data) => {
            if (data.createModel) {
                localStorage.setItem(localStorageKey, JSON.stringify([...models, data.createModel]))
                setModels([...models, data.createModel])
            }
        },
    })

    const addModel = (model: Model) => {
        createModelMutation({
            variables: {
                fields: model.fields,
                name: model.name,
            },
        })
    }

    return (
        <ModelsContext.Provider value={{ models, addModel, sendInformation: handleSendInformation }}>
            {children}
        </ModelsContext.Provider>
    )
}
