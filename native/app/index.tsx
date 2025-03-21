import { client } from "@/api/client";
import Button from "@/components/Button";
import { Alert, Text, View } from "react-native";

export default function Index() {

  const mutateLogin = client.path('/login').method('post').create()

  const handleLogin = () => {
   mutateLogin({
      email: "Kevin",
      password: "Kevin"
    }).then(({ data, status }) => {
      Alert.alert(`Status: ${status}`, JSON.stringify(data))
    }).catch((e) => {
      if (e instanceof mutateLogin.Error) {
        const { data, status } = e.getActualType()
        Alert.alert(`Error: ${status}`, JSON.stringify(data.error))
      } else {
        Alert.alert('Error', JSON.stringify(e))
      }
    })
  }

  return (
    <View className="flex-1 items-center justify-center gap-y-2">
      <View className="items-center">
        <Text className="text-4xl">Starter Template</Text>
        <Text className="text-xl">React Native + Golang</Text>
      </View>
      <Button
        label="Click to Login"
        onPress={handleLogin}
      />
    </View>
  );
}
