//
//  Book.swift
//  restIII
//
//  Created by nicola de filippo on 30/11/23.
//

import Foundation

struct GoogleBooks: Codable {
    var kind: String = ""
    var totalItems: Int = 0
    var items: Array<GoogleBook> = []
}

struct GoogleBook: Codable {
    var id: String
    var volumeInfo: VolumeInfo?
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        id = try container.decode(String.self, forKey: .id)
        volumeInfo = try container.decodeIfPresent(VolumeInfo.self, forKey: .volumeInfo)
    }
}

struct VolumeInfo: Codable {
    var title: String
    var subtitle: String?
    var authors: Array<String>?
    var imageLinks: ImagesLink
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        subtitle = try container.decodeIfPresent(String.self, forKey: .subtitle)
        title = try container.decode(String.self, forKey: .title)
        authors = try container.decodeIfPresent(Array<String>.self, forKey: .authors)
        imageLinks = try container.decode(ImagesLink.self, forKey: .imageLinks)
    }
}

struct ImagesLink: Codable {
    var smallThumbnail: String
}



@Observable
class BookViewModel {
    var books: [GoogleBook] = []
    var state: State = .Loaded
    
    enum State {
        case Loading, Loaded, Error
    }
    
    func loadBooks() async {
        state = .Loading
        guard let url = URL(string: "https://www.googleapis.com/books/v1/volumes?q=intitle:swift") else {
            print("Invalid URL")
            state = .Error
            return
        }
        do {
            let (data, _) = try await URLSession.shared.data(from: url)
            let decoder = JSONDecoder()
            do {
                let decodedResponse = try decoder.decode(GoogleBooks.self, from: data)
                books = decodedResponse.items
                state = .Loaded
                return
            } catch let DecodingError.dataCorrupted(context) {
                print(context)
            } catch let DecodingError.keyNotFound(key, context) {
                print("Key '\(key)' not found:", context.debugDescription)
                print("codingPath:", context.codingPath)
            } catch let DecodingError.valueNotFound(value, context) {
                print("Value '\(value)' not found:", context.debugDescription)
                print("codingPath:", context.codingPath)
            } catch let DecodingError.typeMismatch(type, context)  {
                print("Type '\(type)' mismatch:", context.debugDescription)
                print("codingPath:", context.codingPath)
            } catch {
                print("error: ", error)
            }
        } catch {
            print("error: ", error)
            
        }
        state = .Error
    }
}
