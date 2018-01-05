
# misc notes

Example using the ",any,attr" xml dsl

```go
    package main
    
    import (
    	"encoding/json"
    	"encoding/xml"
    	"fmt"
    )
    
    func main() {
    	type Email struct {
    		XMLName xml.Name
    		Where   string `xml:"where,attr"`
    		Addr    string
    		Attrs   []xml.Attr `xml:",any,attr"`
    	}
    	type Address struct {
    		City, State string
    	}
    	type Result struct {
    		XMLName xml.Name `xml:"Person"`
    		Name    string   `xml:"FullName"`
    		Phone   string
    		Email   []Email
    		Groups  []string `xml:"Group>Value"`
    		Address
    	}
    	v := Result{Name: "none", Phone: "none"}
    	data := `
    		<Person>
    			<FullName>Grace R. Emlin</FullName>
    			<Company>Example Inc.</Company>
    			<Email where="home" preferred="true" is_secret="true">
    				<Addr>gre@example.com</Addr>
    			</Email>
    			<Email where='work' preferred="false">
    				<Addr>gre@work.com</Addr>
    			</Email>
    			<Group>
    				<Value>Friends</Value>
    				<Value>Squash</Value>
    			</Group>
    			<City>Hanga Roa</City>
    			<State>Easter Island</State>
    		</Person>
    	`
    	err := xml.Unmarshal([]byte(data), &v)
    	if err != nil {
    		fmt.Printf("error: %v", err)
    		return
    	}
    	src, _ := json.MarshalIndent(v, "", " ")
    	fmt.Printf("json: %s", src)
    }
```
