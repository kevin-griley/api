import { $api } from '@/api/client';
import Button from '@/components/Button';

import { Alert, Text, View } from 'react-native';

export default function Index() {
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
        <View className="flex-1 items-center justify-center gap-y-2">
            <View className="items-center">
                <Text className="text-4xl">Starter Template</Text>
                <Text className="text-xl">React Native + Golang</Text>
            </View>
            <Button
                label="Click to Login"
                onPress={() =>
                    mutateAsync({
                        body: {
                            email: 'Kevin',
                            password: 'Kevin',
                        },
                    })
                }
            />
        </View>
    );
}
