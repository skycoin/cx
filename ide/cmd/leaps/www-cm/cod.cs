using System;

namespace Example
{
    /// <summary>
    /// Represents a person employed at the company
    /// </summary>
    public class Employee : Person
    {
        #region Properties
        /// <summary>
        /// Gets or sets the first name.
        /// </summary>
        /// <value>The first name.</value>
        public string FirstName { get; set; }

        #endregion

        /// <summary>
        /// Calculates the salary.
        /// </summary>
        /// <param name="grade">The grade.</param>
        /// <returns></returns>
        public decimal CalculateSalary(int grade)
        {
            if (grade > 10)
                return 1000;
            return 500;
        }
    }
}
