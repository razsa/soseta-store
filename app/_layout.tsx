import { Link, Stack } from "expo-router";
import { ShoppingCart } from "lucide-react-native";
import { Icon } from "@/components/ui/icon";

import "@/global.css";
import { GluestackUIProvider } from "@/components/ui/gluestack-ui-provider";
import { Pressable } from "react-native";

export default function RootLayout() {
  return <GluestackUIProvider mode="light">
    <Stack 
    screenOptions={{ headerRight: () => (
      <Link href={'/cart'} asChild>
<Pressable>
  <Icon as={ShoppingCart}/>
</Pressable>
</Link>
),
}}
>
    <Stack.Screen name="index" options={{ title: 'Shop' }} />
    <Stack.Screen name="product/[id]" options={{ title: 'Product' }} />
    </Stack>
    </GluestackUIProvider>;
}
