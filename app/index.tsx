import { View, FlatList } from "react-native";
import products from "../assets/products.json";
import ProductListItem from "../components/ProductListItem";
import { useBreakpointValue } from "@/components/ui/utils/use-break-point-value";
export default function Index() {
  // const { width } = useWindowDimensions();
  // const numColumns = width > 700 ? 3 : 2;
  const numColumns = useBreakpointValue({
    default: 2,
    sm: 3,
  });

  return (
    <View style={{ flex: 1, padding: 10 }}>
      <FlatList
        key={`list-${numColumns}`}
        data={products}
        numColumns={numColumns}
        contentContainerClassName="gap-2 max-w-[960px] mx-auto w-full"
        columnWrapperClassName="gap-2"
        renderItem={({ item }) => <ProductListItem product={item} />}
        keyExtractor={(item) => item.id.toString()}
      />
    </View>
  );
}
