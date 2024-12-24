import { Text } from "@/components/ui/text";
import { useLocalSearchParams, Stack } from "expo-router";
import { usePocketbase, Products } from "@/hooks/pocketbase";
import { Card } from "@/components/ui/card";
import { Image } from "@/components/ui/image";
import { VStack } from "@/components/ui/vstack";
import { Heading } from "@/components/ui/heading";
import { Button, ButtonText } from "@/components/ui/button";
import { Box } from "@/components/ui/box";
import useCart from "@/store/cartStore";
import { useEffect, useState } from "react";
import { getPocketBaseImageUrl } from "@/components/utils/pocketbase-image-url";

export default function ProductDetailsScreen() {
  const [product, setProduct] = useState<Products | null>(null);
  const addProductToCart = useCart((state: CartState) => state.addProduct);
  const { id } = useLocalSearchParams<{ id: string }>();
  const { getProductById } = usePocketbase();

  useEffect(() => {
    const fetchProduct = async () => {
      if (id) {
        const fetchedProduct = await getProductById(id) as Products;
        setProduct(fetchedProduct);
      }
    };

    fetchProduct();
  }, [id]);

  const addToCart = () => {
    if (product) {
      addProductToCart(product);
    }
  };

  if (!product) {
    return <Text>Product not found</Text>;
  }

  return (
    <Box className="bg-purple-200 flex-1 items-center p-3">
      <Stack.Screen options={{ title: product.name }} />
      <Card className="p-5 rounded-lg max-w-[960px] w-full flex-1">
        {product.image && (
          <Image
            source={{
              uri: getPocketBaseImageUrl('products', product.id, product.image),
            }}
            className="mb-6 h-[240px] w-full rounded-md aspect-[4/3]"
            alt={`${product.name} image`}
            resizeMode="contain"
          />
        )}
        <Text className="text-sm font-normal mb-2 text-typography-700">
          {product.name}
        </Text>
        <VStack className="mb-6">
          <Heading size="md" className="mb-4">{product.price}</Heading>
          <Text size="sm">{product.description}</Text>
        </VStack>
        <Box className="flex-col sm:flex-row">
          <Button onPress={addToCart} className="px-4 py-2 mr-0 mb-3 sm:mr-3 sm:mb-0 sm:flex-1">
            <ButtonText size="sm">Add to cart</ButtonText>
          </Button>
          <Button
            variant="outline"
            className="px-4 py-2 border-outline-300 sm:flex-1"
          >
            <ButtonText size="sm" className="text-typography-600">
              Wishlist
            </ButtonText>
          </Button>
        </Box>
      </Card>
    </Box>
  );
}
