//
//  ContentView.swift
//  restIII
//
//  Created by nicola de filippo on 30/11/23.
//

import SwiftUI

struct ContentView: View {
    @State var booksModel = BookViewModel()
    var body: some View {
        NavigationStack {
            switch booksModel.state {
            case .Loading:
                ProgressView()
            case .Error:
                Text("Error")
            case .Loaded:
                List {
                    ForEach(booksModel.books, id:\.id) { book in
                        if let volumeInfo = book.volumeInfo {
                            Text(volumeInfo.title)
                        }
                    }
                }
            }
        }.task {
            await booksModel.loadBooks()
        }
    }
}

#Preview {
    ContentView()
}
