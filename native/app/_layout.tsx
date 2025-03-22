import { Stack } from 'expo-router';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '@/api/client';
import {
    DarkTheme,
    DefaultTheme,
    ThemeProvider,
} from '@react-navigation/native';
import { useColorScheme } from 'react-native';
import { TamaguiProvider } from 'tamagui';
import { tamaguiConfig } from '../tamagui.config';

export default function RootLayout() {
    const colorScheme = useColorScheme();
    return (
        <QueryClientProvider client={queryClient}>
            <TamaguiProvider config={tamaguiConfig} defaultTheme={colorScheme!}>
                <ThemeProvider
                    value={colorScheme === 'dark' ? DarkTheme : DefaultTheme}
                >
                    <Stack />
                </ThemeProvider>
            </TamaguiProvider>
        </QueryClientProvider>
    );
}
