import { $api } from '@/api/client';
import {
    Inter_400Regular,
    Inter_900Black,
    useFonts,
} from '@expo-google-fonts/inter';
import { useEffect } from 'react';
import React, { useState } from 'react';
import { Alert } from 'react-native';
import { Button, Input, Stack, Text, YStack } from 'tamagui';

export default function Page() {
    useFonts({ Inter_400Regular, Inter_900Black });

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const { mutateAsync } = $api.useMutation('post', '/login', {
        onSuccess: (data) => {
            console.log(data);
            Alert.alert('Success', 'Login Successful');
        },
        onError: (error) => {
            console.error(error);
            Alert.alert(`Error: ${error.status}`, error.error);
        },
    });

    return (
        <YStack jc="center" ai="center" p="$4">
            <Stack w={300} p="$4" br="$4">
                <Text fontSize="$8" fontWeight="bold" mb="$3">
                    Welcome
                </Text>
                <Text fontSize="$6" mb="$8">
                    Expo + Golang App
                </Text>
                <Text fontSize="$6" fontWeight="bold" mb="$3">
                    Login
                </Text>

                <Input
                    placeholder="Email"
                    value={email}
                    onChangeText={setEmail}
                    mb="$3"
                />

                <Input
                    placeholder="Password"
                    value={password}
                    onChangeText={setPassword}
                    secureTextEntry
                    mb="$3"
                />

                <Button
                    onPress={() =>
                        mutateAsync({
                            body: {
                                email,
                                password,
                            },
                        })
                    }
                    bg="$color6"
                >
                    Login
                </Button>
            </Stack>
        </YStack>
    );
}
