import createFetchClient from 'openapi-fetch';
import createClient from 'openapi-react-query';
import type { paths } from '@/types/schema';
import { QueryClient } from '@tanstack/react-query';

const fetchClient = createFetchClient<paths>({
    baseUrl: 'http://localhost:3000',
    headers: {
        'Content-Type': 'application/json',
    },
});

export const $api = createClient(fetchClient);

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            retry: false,
        },
    },
});
