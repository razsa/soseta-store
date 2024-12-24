import { View, FlatList, useWindowDimensions } from "react-native";
import ProductListItem from "../components/ProductListItem";
import { useBreakpointValue } from "../components/ui/utils/use-break-point-value";
import { usePocketbase } from "../hooks/pocketbase";
import { useEffect, useState } from "react";

interface Product {
  id: string;
  name: string;
  price: number;
  image: string;
  collectionId: string;
  description: string;
  stock: string;
  collectionName: string;
  imageUrl: string;
}

export default function Index() {
  const { getProducts } = usePocketbase();
  const [data, setData] = useState<Product[]>([]);

  useEffect(() => {
    const loadData = async () => {
      const result = await getProducts();
      console.log('Product data:', result);
      setData(result);
    };
    loadData();
  }, []);

  const { width } = useWindowDimensions();
  
  const numColumns = useBreakpointValue({
    default: 2,     // Mobile (<640px)
    md: 3,          // Tablet (≥768px)
    lg: 4,          // Desktop (≥1024px)
  });

  return (
    <View style={{ flex: 1 }}>
      <FlatList
        key={`list-${numColumns}`}
        data={data}
        numColumns={numColumns}
        contentContainerStyle={{ padding: 8 }}
        columnWrapperStyle={{ 
          gap: 8,
          justifyContent: 'flex-start',
          paddingHorizontal: 8
        }}
        ItemSeparatorComponent={() => <View style={{ height: 8 }} />}
        renderItem={({ item }) => <ProductListItem product={item} />}
        keyExtractor={(item) => item.id.toString()}
      />
    </View>
  );
}
