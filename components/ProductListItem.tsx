import { Text, Pressable, useWindowDimensions, View } from "react-native";
import { Card } from "@/components/ui/card";
import { Image } from "@/components/ui/image";
import { Heading } from "@/components/ui/heading";
import { Link } from "expo-router";
import { useBreakpointValue } from "./ui/utils/use-break-point-value";
import { getPocketBaseImageUrl } from "./utils/pocketbase-image-url";

interface Product {
  id: string;
  name: string;
  price: number;
  image: string;
  collectionId: string;
  description: string;
  stock: string;
  collectionName: string;
}

export default function ProductListItem({ product }: { product: Product }) {
  const { width: screenWidth } = useWindowDimensions();
  
  const cardWidth = useBreakpointValue({
    default: Math.floor((screenWidth - 32) / 2),
    md: 230,
    lg: 250,
  });

  // Construct the PocketBase image URL
  const imageUrl = getPocketBaseImageUrl(product.collectionId, product.id, product.image);
  
  console.log('Constructed image URL:', imageUrl); // Debugging log

  return (
    <Link href={`/product/${product.id}`} asChild>
      <Pressable style={{ width: cardWidth }}>
        <Card className="p-4 rounded-lg">
          {product.image && (
              <Image
                source={{ uri: imageUrl }}
                className="mb-4 h-[160px] w-full rounded-md aspect-[4/3]"
                resizeMode="contain"
                alt={`${product.name} image`}
              />
            )}
          
          <Text className="text-sm font-normal mt-4 mb-2 text-typography-700" numberOfLines={1}>
            {product.name}
          </Text>
          <Heading size="md" className="mb-2">
            ${product.price.toFixed(2)}
          </Heading>
        </Card>
      </Pressable>
    </Link>
  );
}
