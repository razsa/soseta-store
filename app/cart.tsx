import React from "react";
import { FlatList, Text, View } from "react-native";
import { VStack } from "../components/ui/vstack";
import { Button, ButtonText } from "../components/ui/button";
import useCart from "@/store/cartStore";
import { Redirect } from "expo-router";

const onCheckout = () => {
  console.log("Checkout pressed");
};

export default function CartScreen() {
  const items = useCart((state) => state.items);
  const resetCart = useCart((state) => state.resetCart);

  const onCheckout = async () => {
    // send order to server
    resetCart();
  };

  if (items.length === 0) {
    return <Redirect href={"/"} />
  }

  return (
    <FlatList
      data={items}
      contentContainerClassName="gap-2 max-w-[960px] w-full mx-auto p-2"
      renderItem={({ item }) => (
        <VStack className="bg-white p-3">
          <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' }}>
            <VStack space="sm">
              <Text style={{ fontWeight: 'bold' }}>{item.product.name}</Text>
              <Text>$ {item.product.price}</Text>
            </VStack>
            <Text className="ml-auto">{item.quantity}</Text>
          </View>
        </VStack>
      )}
      ListFooterComponent={() => (
        <Button onPress={onCheckout}>
          <ButtonText>Checkout</ButtonText>
        </Button>
      )}
    />
  );
}
