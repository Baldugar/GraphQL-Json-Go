import { Button, Container, Grid, Typography } from '@mui/material'
import MDEditor from '@uiw/react-md-editor'
import { useContext, useState } from 'react'
import { ModelsContext } from '../../context/ModelsContext'

const Home = () => {
    const { sendInformation } = useContext(ModelsContext)
    const [info, setInfo] = useState<string>('')

    return (
        <Container>
            <Typography variant={'h3'}>Models</Typography>
            <Grid container spacing={3}>
                <Grid item xs={12}>
                    <MDEditor value={info} onChange={(v) => setInfo(v || '')} />
                </Grid>
                <Grid item xs={12}>
                    <Button
                        onClick={() => {
                            const trimmed = info.trim().replace(/\n/g, '')
                            if (trimmed.length > 0) {
                                sendInformation(JSON.parse(trimmed), 'Person')
                            }
                        }}
                        variant={'contained'}
                    >
                        Add Model
                    </Button>
                </Grid>
            </Grid>
        </Container>
    )
}

export default Home
