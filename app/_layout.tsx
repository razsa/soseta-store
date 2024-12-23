import { Link, Stack } from "expo-router";
import { ShoppingCart, User } from "lucide-react-native";
import { Icon } from "@/components/ui/icon";
import useCart from "@/store/cartStore";
import { GluestackUIProvider } from "@/components/ui/gluestack-ui-provider";
import { Pressable, Text } from "react-native";

import "@/global.css";

export default function RootLayout() {
  const cart = useCart();
  const itemCount = cart.items.reduce((total, item) => total + item.quantity, 0);

  return (
    <GluestackUIProvider mode="light">
      <Stack
        screenOptions={{
          headerRight: () => (
            <Link href="/cart" asChild>
              <Pressable className="flex-row gap-2">
                <Icon as={ShoppingCart} />
                {itemCount > 0 && (
                  <Text style={{ marginLeft: 5, fontSize: 12 }}>{itemCount}</Text>
                )}
              </Pressable>
            </Link>
          ),
          headerLeft: () => (
            <Link href="/login" asChild>
              <Pressable className="flex-row gap-2">
                <Icon as={User} />
              </Pressable>
            </Link>
          ),
        }}
        >
        <Stack.Screen name="index" options={{ title: "Shop" }} />
        <Stack.Screen name="product/[id]" options={{ title: "Product" }} />
      </Stack>
    </GluestackUIProvider>
  );
}
