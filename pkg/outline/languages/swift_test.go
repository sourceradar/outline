package languages

import (
	"strings"
	"testing"

	swift "github.com/alex-pinkus/tree-sitter-swift/bindings/go"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

func TestSwiftOutlineWithImports(t *testing.T) {
	swiftCode := `import UIKit
import Foundation

/// A simple class for demonstration
public class MyViewController: UIViewController {
    
    /// The main view
    @IBOutlet weak var mainView: UIView!
    
    /// Initialize the controller
    public override init(nibName: String?, bundle: Bundle?) {
        super.init(nibName: nibName, bundle: bundle)
    }
    
    /// Required initializer
    required init?(coder: NSCoder) {
        super.init(coder: coder)
    }
    
    /// View did load
    public override func viewDidLoad() {
        super.viewDidLoad()
        setupUI()
    }
    
    /// Setup the user interface
    private func setupUI() {
        // Implementation here
    }
    
    /// Cleanup resources
    deinit {
        // Cleanup code
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that imports are included
	if !strings.Contains(result, "import UIKit") {
		t.Error("Expected import UIKit to be included")
	}
	if !strings.Contains(result, "import Foundation") {
		t.Error("Expected import Foundation to be included")
	}

	// Check that class is included with modifiers and inheritance
	if !strings.Contains(result, "public class MyViewController: UIViewController") {
		t.Error("Expected class declaration with modifiers and inheritance to be included")
	}

	// Check that properties are included
	if !strings.Contains(result, "mainView: UIView") {
		t.Error("Expected property declaration to be included")
	}

	// Check that methods are included
	if !strings.Contains(result, "func viewDidLoad()") {
		t.Error("Expected method declaration to be included")
	}

	// Check that init methods are included
	if !strings.Contains(result, "init(nibName: String?, bundle: Bundle?)") {
		t.Error("Expected init method to be included")
	}

	// Check that deinit is included
	if !strings.Contains(result, "deinit") {
		t.Error("Expected deinit to be included")
	}
}

func TestSwiftStructWithProtocols(t *testing.T) {
	swiftCode := `
/// A point in 2D space
public struct Point: Codable, Equatable {
    /// X coordinate
    let x: Double
    /// Y coordinate  
    let y: Double
    
    /// Initialize with coordinates
    public init(x: Double, y: Double) {
        self.x = x
        self.y = y
    }
    
    /// Calculate distance to another point
    public func distance(to other: Point) -> Double {
        return sqrt(pow(x - other.x, 2) + pow(y - other.y, 2))
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that struct is included with modifiers and protocols
	if !strings.Contains(result, "public struct Point: Codable, Equatable") {
		t.Error("Expected struct declaration with protocols to be included")
	}

	// Check that properties are included
	if !strings.Contains(result, "x: Double") {
		t.Error("Expected x property to be included")
	}
	if !strings.Contains(result, "y: Double") {
		t.Error("Expected y property to be included")
	}

	// Check that methods are included
	if !strings.Contains(result, "func distance(to other: Point) -> Double") {
		t.Error("Expected method with parameters and return type to be included")
	}
}

func TestSwiftProtocol(t *testing.T) {
	swiftCode := `
/// A drawable protocol
public protocol Drawable {
    /// Draw the object
    func draw()
    
    /// The drawing color
    var color: UIColor { get set }
    
    /// Optional method
    func animate() -> Void
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that protocol is included
	if !strings.Contains(result, "public protocol Drawable") {
		t.Error("Expected protocol declaration to be included")
	}

	// Check that protocol methods are included
	if !strings.Contains(result, "func draw()") {
		t.Error("Expected protocol method to be included")
	}

	// Check that protocol properties are included
	if !strings.Contains(result, "color: UIColor") {
		t.Error("Expected protocol property to be included")
	}
}

func TestSwiftEnum(t *testing.T) {
	swiftCode := `
/// HTTP status codes
public enum HTTPStatus: Int {
    case ok = 200
    case notFound = 404
    case serverError = 500
    
    /// Get status message
    public func message() -> String {
        switch self {
        case .ok:
            return "OK"
        case .notFound:
            return "Not Found"
        case .serverError:
            return "Server Error"
        }
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that enum is included with raw type
	if !strings.Contains(result, "public enum HTTPStatus: Int") {
		t.Error("Expected enum declaration with raw type to be included")
	}

	// Check that enum cases are included
	if !strings.Contains(result, "case ok, notFound, serverError") {
		t.Error("Expected enum cases to be included")
	}

	// Check that enum methods are included
	if !strings.Contains(result, "func message() -> String") {
		t.Error("Expected enum method to be included")
	}
}

func TestSwiftExtension(t *testing.T) {
	swiftCode := `
/// String extensions
extension String: CustomStringConvertible {
    
    /// Check if string is empty or whitespace
    public var isBlank: Bool {
        return trimmingCharacters(in: .whitespacesAndNewlines).isEmpty
    }
    
    /// Capitalize first letter
    public func capitalizeFirst() -> String {
        return prefix(1).uppercased() + dropFirst()
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that extension is included with protocol conformance
	if !strings.Contains(result, "extension String: CustomStringConvertible") {
		t.Error("Expected extension with protocol conformance to be included")
	}

	// Check that computed properties are included
	if !strings.Contains(result, "isBlank: Bool") {
		t.Error("Expected computed property to be included")
	}

	// Check that extension methods are included
	if !strings.Contains(result, "func capitalizeFirst() -> String") {
		t.Error("Expected extension method to be included")
	}
}

func TestSwiftTypealias(t *testing.T) {
	swiftCode := `
/// Completion handler
public typealias CompletionHandler = (Bool) -> Void

/// Point type alias
typealias Point2D = (x: Double, y: Double)
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that typealiases are included
	if !strings.Contains(result, "public typealias CompletionHandler = (Bool) -> Void") {
		t.Error("Expected public typealias to be included")
	}

	if !strings.Contains(result, "typealias Point2D = (x: Double, y: Double)") {
		t.Error("Expected typealias to be included")
	}
}

func TestSwiftSubscript(t *testing.T) {
	swiftCode := `
public struct Matrix {
    private var data: [[Double]]
    
    /// Access matrix elements
    public subscript(row: Int, column: Int) -> Double {
        get {
            return data[row][column]
        }
        set {
            data[row][column] = newValue
        }
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that struct is included
	if !strings.Contains(result, "public struct Matrix") {
		t.Error("Expected struct declaration to be included")
	}

	// Check that subscript is included
	if !strings.Contains(result, "subscript(row: Int, column: Int) -> Double") {
		t.Error("Expected subscript declaration to be included")
	}
}

func TestSwiftComplexClass(t *testing.T) {
	swiftCode := `
import Foundation

/// A network manager class
@objc public class NetworkManager: NSObject {
    
    /// Shared instance
    public static let shared = NetworkManager()
    
    /// Base URL
    private let baseURL: URL
    
    /// Session configuration
    private lazy var session: URLSession = {
        let config = URLSessionConfiguration.default
        return URLSession(configuration: config)
    }()
    
    /// Private initializer
    private override init() {
        self.baseURL = URL(string: "https://api.example.com")!
        super.init()
    }
    
    /// Make a GET request
    public func get<T: Codable>(
        endpoint: String,
        type: T.Type,
        completion: @escaping (Result<T, Error>) -> Void
    ) {
        // Implementation
    }
    
    /// Download data
    public func download(
        from url: URL,
        progress: ((Double) -> Void)? = nil,
        completion: @escaping (Result<Data, Error>) -> Void
    ) {
        // Implementation
    }
}
`

	parser := sitter.NewParser()
	defer parser.Close()

	if err := parser.SetLanguage(sitter.NewLanguage(swift.Language())); err != nil {
		t.Fatalf("Failed to set Swift language: %v", err)
	}

	tree := parser.Parse([]byte(swiftCode), nil)
	defer tree.Close()

	result := ExtractSwiftOutline(tree.RootNode(), []byte(swiftCode))

	// Check that import is included
	if !strings.Contains(result, "import Foundation") {
		t.Error("Expected import to be included")
	}

	// Check that class with attributes is included
	if !strings.Contains(result, "public class NetworkManager: NSObject") {
		t.Error("Expected class declaration to be included")
	}

	// Check that static properties are included
	if !strings.Contains(result, "shared") {
		t.Error("Expected static property to be included")
	}

	// Check that complex methods with generics and closures are included
	if !strings.Contains(result, "func get") {
		t.Error("Expected generic method to be included")
	}

	if !strings.Contains(result, "func download") {
		t.Error("Expected method with optional parameters to be included")
	}
}
