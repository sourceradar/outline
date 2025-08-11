package languages

import (
	"strings"
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	c "github.com/tree-sitter/tree-sitter-c/bindings/go"
	cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
)

func TestCOutlineWithDefinesAndFunctions(t *testing.T) {
	cCode := `#include <stdio.h>
#include <stdlib.h>

#define MAX_SIZE 100
#define SQUARE(x) ((x) * (x))

/* Global variable */
int global_var = 42;

/* 
 * Main function
 * Returns 0 on success
 */
int main(int argc, char *argv[]) {
    printf("Hello, World!\n");
    return 0;
}

// Helper function
void helper_function(int value) {
    printf("Value: %d\n", value);
}

/* Function pointer declaration */
int (*func_ptr)(int, int);
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(c.Language())); err != nil {
		t.Fatalf("Failed to set C language: %v", err)
	}

	tree := parser.Parse([]byte(cCode), nil)
	defer tree.Close()

	result := ExtractCOutline(tree.RootNode(), []byte(cCode))

	// Check that includes are included
	if !strings.Contains(result, "#include <stdio.h>") {
		t.Error("Expected #include <stdio.h> to be included")
	}

	// Check that defines are included
	if !strings.Contains(result, "#define MAX_SIZE 100") {
		t.Error("Expected #define MAX_SIZE to be included")
	}

	if !strings.Contains(result, "#define SQUARE(x)") {
		t.Error("Expected #define SQUARE(x) to be included")
	}

	// Check that functions are included
	if !strings.Contains(result, "int main(int argc, char *argv[])") {
		t.Error("Expected main function to be included")
	}

	if !strings.Contains(result, "void helper_function(int value)") {
		t.Error("Expected helper_function to be included")
	}

	// Check that comments are included
	if !strings.Contains(result, "Main function") {
		t.Error("Expected function documentation to be included")
	}
}

func TestCOutlineWithStructsAndEnums(t *testing.T) {
	cCode := `#include <stdint.h>

/* Point structure */
struct Point {
    int x;
    int y;
};

/* Union for different data types */
union Data {
    int integer;
    float floating;
    char string[20];
};

/* Color enumeration */
enum Color {
    RED = 1,
    GREEN,
    BLUE
};

typedef struct {
    char name[50];
    int age;
} Person;

typedef enum {
    STATUS_OK,
    STATUS_ERROR,
    STATUS_PENDING
} Status;

/* Function using struct */
struct Point create_point(int x, int y) {
    struct Point p = {x, y};
    return p;
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(c.Language())); err != nil {
		t.Fatalf("Failed to set C language: %v", err)
	}

	tree := parser.Parse([]byte(cCode), nil)
	defer tree.Close()

	result := ExtractCOutline(tree.RootNode(), []byte(cCode))

	// Check that struct is included
	if !strings.Contains(result, "struct Point {") {
		t.Error("Expected struct Point to be included")
	}

	// Check that union is included
	if !strings.Contains(result, "union Data {") {
		t.Error("Expected union Data to be included")
	}

	// Check that enum is included
	if !strings.Contains(result, "enum Color {") {
		t.Error("Expected enum Color to be included")
	}

	// Check that typedef is included
	if !strings.Contains(result, "typedef struct") {
		t.Error("Expected typedef struct to be included")
	}

	// Check that function using struct is included
	if !strings.Contains(result, "struct Point create_point") {
		t.Error("Expected create_point function to be included")
	}
}

func TestCppOutlineWithClassesAndNamespaces(t *testing.T) {
	cppCode := `#include <iostream>
#include <vector>
#include <string>

namespace Math {
    /* Mathematical constants */
    const double PI = 3.14159;
    
    class Calculator {
    private:
        std::string name;
        
    public:
        Calculator(const std::string& n) : name(n) {}
        
        // Basic arithmetic operations
        double add(double a, double b) {
            return a + b;
        }
        
        double multiply(double a, double b);
        
        virtual ~Calculator() = default;
    };
    
    class ScientificCalculator : public Calculator {
    private:
        bool degree_mode;
        
    public:
        ScientificCalculator() : Calculator("Scientific") {}
        
        double sin(double angle);
        double cos(double angle);
    };
}

template<typename T>
class Vector {
private:
    std::vector<T> data;
    
public:
    void push(const T& item) {
        data.push_back(item);
    }
    
    T get(size_t index) const {
        return data[index];
    }
};

int main() {
    Math::Calculator calc("Basic");
    return 0;
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(cpp.Language())); err != nil {
		t.Fatalf("Failed to set C++ language: %v", err)
	}

	tree := parser.Parse([]byte(cppCode), nil)
	defer tree.Close()

	result := ExtractCppOutline(tree.RootNode(), []byte(cppCode))

	expected := `#include <iostream>

#include <vector>

#include <string>

namespace Math { // line 5
	const double PI = 3.14159; // line 7
	class Calculator { // line 9
	private:
	public:
		Calculator(const std::string& n) { //... } // line 14
		double add(double a, double b) { //... } // line 17
		~Calculator() { //... } // line 23
	}

	class ScientificCalculator : : public Calculator { // line 26
	private:
	public:
		ScientificCalculator() { //... } // line 31
	}

}

template<typename T> // line 38
class Vector { // line 39
	private:
	public:
		void push(const T& item) { //... } // line 44
		T get(size_t index) const { //... } // line 48
}

int main() { //... } // line 53

`

	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestCOutlineWithComplexStructures(t *testing.T) {
	cCode := `/* Network packet structure */
struct packet {
    uint16_t length;
    uint8_t type;
    union {
        struct {
            uint32_t src_addr;
            uint32_t dest_addr;
        } ip;
        struct {
            uint16_t port;
            uint8_t flags;
        } tcp;
    } header;
    uint8_t data[1024];
};

/* Function prototypes */
int send_packet(struct packet *pkt);
int receive_packet(struct packet *pkt, int timeout);
void process_packet(const struct packet *pkt);

/* Callback function type */
typedef int (*packet_handler_t)(struct packet *pkt);

/* Error codes */
enum net_error {
    NET_OK = 0,
    NET_TIMEOUT = -1,
    NET_INVALID_PACKET = -2,
    NET_BUFFER_FULL = -3
};
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(c.Language())); err != nil {
		t.Fatalf("Failed to set C language: %v", err)
	}

	tree := parser.Parse([]byte(cCode), nil)
	defer tree.Close()

	result := ExtractCOutline(tree.RootNode(), []byte(cCode))

	// Check that complex struct is included
	if !strings.Contains(result, "struct packet {") {
		t.Error("Expected struct packet to be included")
	}

	// Check that function prototypes are included
	if !strings.Contains(result, "int send_packet(struct packet *pkt);") {
		t.Error("Expected send_packet prototype to be included")
	}

	// Check that typedef for function pointer is included
	if !strings.Contains(result, "typedef int (*packet_handler_t)") {
		t.Error("Expected packet_handler_t typedef to be included")
	}

	// Check that enum with values is included
	if !strings.Contains(result, "enum net_error {") {
		t.Error("Expected enum net_error to be included")
	}
}
