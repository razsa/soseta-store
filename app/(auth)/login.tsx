import { FormControl } from "@/components/ui/form-control";
import { Text } from "@/components/ui/text";
import { Input, InputField, InputIcon, InputSlot } from "@/components/ui/input";
import { HStack } from "@/components/ui/hstack";
import { Heading } from '@/components/ui/heading';
import { VStack } from "@/components/ui/vstack";
import { EyeIcon, EyeOffIcon } from "lucide-react-native";
import { Button, ButtonText } from "@/components/ui/button";
import { useState } from "react";

export default function LoginScreen() {
    const [showPassword, setShowPassword] = useState(false);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const handleState = () => {
        setShowPassword((showState) => {
            return !showState;
        });
    };

    const signIn = async () => {
        try {
            const response = await fetch('http://localhost:8090/api/collections/users/auth-with-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    identity: email,
                    password: password,
                }),
            });

            if (response.ok) {
                window.location.href = '/';
            } else {
                console.error('Login failed');
                const text = await response.text();
                console.log(text);
            }
        } catch (error) {
            console.error('Error during login:', error);
        }
    };

    const signUp = async () => {
        try {
            const response = await fetch('http://localhost:8090/api/collections/users/records', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email: email,
                    password: password,
                    passwordConfirm: password,
                }),
            });

            if (response.ok) {
                // Redirect to home page on successful signup
                window.location.href = '/';
            } else {
                console.error('Signup failed');
                const text = await response.text();
                console.log(text);
                // Optionally display an error message to the user
            }
        } catch (error) {
            console.error('Error during signup:', error);
            // Optionally display an error message to the user
        }
    };

    return (
        <FormControl
            className="p-4 border rounded-lg max-w-[500px] border-outline-300 bg-white m-2"
        >
            <VStack space="xl">
                <Heading className="text-typography-900 leading-3 pt-3">Login</Heading>
                <VStack space="xs">
                    <Text className="text-typography-500 leading-1">Email</Text>
                    <Input>
                        <InputField type="text" value={email} onChangeText={setEmail} />
                    </Input>
                </VStack>
                <VStack space="xs">
                    <Text className="text-typography-500 leading-1">Password</Text>
                    <Input className="text-center">
                        <InputField
                            type={showPassword ? 'text' : 'password'}
                            value={password}
                            onChangeText={setPassword}
                        />
                        <InputSlot className="pr-3" onPress={handleState}>
                            <InputIcon
                                as={showPassword ? EyeIcon : EyeOffIcon}
                                className="text-darkBlue-500"
                            />
                        </InputSlot>
                    </Input>
                </VStack>
                <HStack space="sm">
                    <Button
                        className="flex-1"
                        variant="outline"
                        onPress={signUp}>
                        <ButtonText>Sign up</ButtonText>
                    </Button>
                    <Button
                        className="flex-1"
                        onPress={signIn}>
                        <ButtonText>Sign in</ButtonText>
                    </Button>
                </HStack>
            </VStack>
        </FormControl>
    );
}
