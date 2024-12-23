import { Text, Pressable, useWindowDimensions } from "react-native";
import { Card } from "@/components/ui/card";
import { Image } from "@/components/ui/image";
import { Heading } from "@/components/ui/heading";
import { Link } from "expo-router";
import { useBreakpointValue } from "./ui/utils/use-break-point-value";

export default function ProductListItem({ product }: { product: Product }) {
    const { width: screenWidth } = useWindowDimensions();
  
  const cardWidth = useBreakpointValue({
    default: Math.floor((screenWidth - 32) / 2), // Mobile: (screen width - total padding) / 2 cards
    md: 230,       // Tablet
    lg: 250,       // Desktop
  });
  return (
    <Link href={`/product/${product.id}`} asChild>
    <Pressable style={{ width: cardWidth}}>
    <Card className="p-4 rounded-lg">
    <Image
      source={{
        uri: product.image,
      }}
      className="mb-4 h-[160px] w-full rounded-md aspect-[4/3]"
      alt={`${product.name} image`}
      resizeMode="contain"
    />
    <Text className="text-sm font-normal mb-2 text-typography-700" numberOfLines={1}>
      {product.name} 
    </Text>
      <Heading size="md" className="mb-2">
       {product.price} 
      </Heading>
  </Card>
  </Pressable>
  </Link>
  );
}
