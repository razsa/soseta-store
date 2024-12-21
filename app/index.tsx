import { Text, View, FlatList } from "react-native";
import products from "../assets/products.json";
import ProductListItem from "../components/ProductListItem";
import { Button, ButtonText } from "@/components/ui/button";

export default function Index() {
  return (
    <View style={{ flex: 1, padding: 10 }}>
      <FlatList
        data={products}
        numColumns={2}
        contentContainerClassName="gap-2"
        columnWrapperClassName="gap-2"
        renderItem={({ item }) => <ProductListItem product={item} />}
        keyExtractor={(item) => item.id}
      />
      <Button>
        <ButtonText>Click me</ButtonText>
      </Button>
    </View>
  );
}
