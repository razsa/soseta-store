import { Text, Pressable } from "react-native";
import { Box } from "@/components/ui/box";
import { Card } from "@/components/ui/card";
import { Image } from "@/components/ui/image";
import { VStack } from "@/components/ui/vstack";
import { Heading } from "@/components/ui/heading";
import { Link } from "expo-router";

export default function ProductListItem({ product }: { product: Product }) {
  return (
    <Link href={`/product/${product.id}`}>
    <Pressable className="flex-1">
    <Card className="p-5 rounded-lg max-w-[360px] flex-1">
    <Image
      source={{
        uri: product.image,
      }}
      className="mb-6 h-[240px] w-full rounded-md aspect-[4/3]"
      alt={`${product.name} image`}
      resizeMode="contain"
    />
    <Text className="text-sm font-normal mb-2 text-typography-700">
      {product.name} 
    </Text>
      <Heading size="md" className="mb-4">
       {product.price} 
      </Heading>
  </Card>
  </Pressable>
  </Link>
  );
}
