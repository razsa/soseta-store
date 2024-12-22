import { View, FlatList, useWindowDimensions } from "react-native";
import products from "../assets/products.json";
import ProductListItem from "../components/ProductListItem";
import { useBreakpointValue } from "@/components/ui/utils/use-break-point-value";

export default function Index() {
  const { width } = useWindowDimensions();
  console.log('Current window width:', width);
  
  const numColumns = useBreakpointValue({
    default: 2,     // Mobile (<640px)
    md: 3,          // Tablet (≥768px)
    lg: 4,          // Desktop (≥1024px)
  });
  console.log('Current numColumns:', numColumns);

  return (
      <FlatList
        key={`list-${numColumns}`}
        data={products}
        numColumns={numColumns}
        contentContainerClassName="gap-2 py-4 px-2"
        columnWrapperClassName="gap-2 justify-center w-full"
        renderItem={({ item }) => <ProductListItem product={item} />}
        keyExtractor={(item) => item.id.toString()}
      />
  );
}
