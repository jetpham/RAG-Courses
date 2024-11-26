from flask import Flask, request, jsonify
import ratemyprofessor
import json

app = Flask(__name__)

# Get the professor's information using the RateMyProfessor API that I could only find in python
def get_professor_info(professor_name):
    school = ratemyprofessor.get_school_by_name("University of San Francisco")
    
    if school is None:
        return {"error": "School not found"}

    # Get the professor
    professor = ratemyprofessor.get_professor_by_school_and_name(school, professor_name)

    if professor is None:
        return {"error": "Professor not found"}

    # Extract and format the professor's information
    professor_info = {
        "name": professor.name,
        "department": professor.department,
        "school": professor.school.name,
        "rating": professor.rating,
        "difficulty": professor.difficulty,
        "total_ratings": professor.num_ratings,
        "would_take_again": round(professor.would_take_again, 1) if professor.would_take_again is not None else "N/A"
    }

    return professor_info

@app.route('/professor', methods=['GET'])
def get_professor():
    professor_name = request.args.get('name')
    if professor_name is None:
        return jsonify({"error": "Professor name is required"}), 400
    
    professor_info = get_professor_info(professor_name)
    return jsonify(professor_info)

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)