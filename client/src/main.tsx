import * as React from 'react'
import { ChakraProvider } from '@chakra-ui/react'
import * as ReactDOM from 'react-dom/client'
import App from './App'
import './index.css';
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import theme from "./chakra/theme.ts";

const queryClient = new QueryClient();
const rootElement = document.getElementById('root')
ReactDOM.createRoot(rootElement!).render(
	<React.StrictMode>
		<QueryClientProvider client={queryClient}>
			<ChakraProvider theme={theme}>
				<App />
			</ChakraProvider>
		</QueryClientProvider>
	</React.StrictMode>
)