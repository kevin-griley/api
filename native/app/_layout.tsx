import { Stack } from 'expo-router';

// Import your global CSS file
import '../global.css';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '@/api/client';

export default function RootLayout() {
    return (
        <QueryClientProvider client={queryClient}>
            <Stack />
        </QueryClientProvider>
    );
}
